// src/presentation/pages/dashboard/dashboard.tsx
import React, {
  createContext,
  lazy,
  ReactElement,
  Suspense,
  useContext,
  useEffect,
  useState,
  useMemo,
  useCallback,
  Component,
  ErrorInfo,
  ReactNode
} from 'react';
import { useTranslation } from 'react-i18next';
import PageListingLoading from '../../components/pages/PageListingLoading';
import { AuthContext } from '../middleware';
import LoadingState from '../../components/pages/loading/loadingState.component';
import { DashboardHeader } from '../../components/header/dashboard-header';
import { UserProvider, UserProviderHook } from '../../providers/user.provider';
import FpsBoostPage from '@/presentation/features/boosters/fps-booster.page';
import ConnectionPage from '@/presentation/features/boosters/connection.page';
import PrecisionPage from '@/presentation/features/boosters/precision.page';
import GamesPage from '@/presentation/features/boosters/games.page';
import FlusherPage from '@/presentation/features/boosters/flusher.page';
import MultiAI from '@/presentation/features/chat-ai/page';
import SmartBoost from '@/presentation/features/mart-booster/smart-booster.page';
import ProfilePage from '@/presentation/features/settings/profile.page';

const ErrorFallback: React.FC<{ messageKey?: string; values?: Record<string, any> }> = ({ messageKey = 'error.homeLoadError', values }) => {
  const { t } = useTranslation();
  return <div className="text-center text-red-400">{t(messageKey, values)}</div>;
};

const HomePage = lazy(() =>
  import('@/presentation/features/homepage/home.page').catch(() => ({
    default: () => <ErrorFallback messageKey="error.homeLoadError" />
  }))
);

export enum PageType {
  HOMEPAGE = "hub",

  FPS_BOOST = "fpsboost",
  CONNECTION = "connection",
  PRECISION = "precision",
  GAMES = "games",
  FLUSHER = "flusher",

  CHAT = "chat",
  ANALYTICS = "analytics",
  SMART_BOOST = "smartboost",

  PROFILE = "profile"
}

// Interface melhorada com validação
interface PagesProviderI {
  currentPage: PageType;
  handleChangePage: (newPage: PageType) => void;
  isLoading: boolean;
  error: string | null;
}

const PAGES_PROVIDER_INITIAL: PagesProviderI = {
  currentPage: PageType.HOMEPAGE,
  handleChangePage: () => {},
  isLoading: false,
  error: null
};

// Context com validação
export const PagesProvider = createContext<PagesProviderI>(PAGES_PROVIDER_INITIAL);

// Error Boundary Component
interface ErrorBoundaryState {
  hasError: boolean;
  error?: Error;
}

class DashboardErrorBoundary extends Component<
  { children: ReactNode; fallback?: ReactElement },
  ErrorBoundaryState
> {
  constructor(props: { children: ReactNode; fallback?: ReactElement }) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Dashboard Error:', error, errorInfo);
    // Aqui você pode enviar o erro para um serviço de monitoramento
  }

  render() {
    const { t } = (this.props as any).__translationProps || {}; // fallback (não usado normalmente)
    if (this.state.hasError) {
      return this.props.fallback || (
        <div className="flex items-center justify-center w-full h-full bg-zinc-950 text-white">
          <div className="text-center">
            <h2 className="text-xl font-semibold mb-2">Algo deu errado</h2>
            <p className="text-gray-400 mb-4">Ocorreu um erro inesperado</p>
            <button
              onClick={() => window.location.reload()}
              className="bg-blue-600 hover:bg-blue-700 px-4 py-2 rounded"
            >
              Recarregar Página
            </button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

// Hook personalizado para gerenciamento de páginas
const usePageManager = () => {
  const [currentPage, setCurrentPage] = useState<PageType>(PageType.HOMEPAGE);
  const [isLoading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const handleChangePage = useCallback((newPage: PageType) => {
    // Validação de entrada
    if (!Object.values(PageType).includes(newPage)) {
      setError('Página inválida solicitada');
      return;
    }

    setError(null);
    setCurrentPage(newPage);
  }, []);

  return {
    currentPage,
    handleChangePage,
    isLoading,
    error
  };
};

// Hook para validação de autenticação
const useAuthValidation = () => {
  const authContext = useContext(AuthContext);
  const [authError, setAuthError] = useState<string | null>(null);

  useEffect(() => {
    if (!authContext) {
      setAuthError('Contexto de autenticação não encontrado');
      return;
    }

    if (authContext.isAuthenticated === false) {
      setAuthError('Usuário não autenticado');
    } else {
      setAuthError(null);
    }
  }, [authContext?.isAuthenticated]);

  return {
    isAuthenticated: authContext?.isAuthenticated ?? false,
    authError
  };
};

// Componente de página com proteção
const PageWrapper: React.FC<{
  page: PageType;
  currentPage: PageType;
  children: ReactElement;
}> = ({ page, currentPage, children }) => {
  const { t } = useTranslation();
  if (page !== currentPage) return null;

  return (
    <DashboardErrorBoundary
      fallback={
        <div className="flex items-center justify-center w-full h-full">
          <p className="text-red-400">{t('error.loadPageFailed', { page })}</p>
        </div>
      }
    >
      <Suspense fallback={<PageListingLoading />}>{children}</Suspense>
    </DashboardErrorBoundary>
  );
};

// Componente principal melhorado
export function DashboardPages() {
  const { t } = useTranslation();
  const { isAuthenticated, authError } = useAuthValidation();
  const pageManager = usePageManager();
  const userProviderValues = UserProviderHook();

  // Memoização dos valores do contexto
  const pagesContextValue = useMemo(
    () => ({
      currentPage: pageManager.currentPage,
      handleChangePage: pageManager.handleChangePage,
      isLoading: pageManager.isLoading,
      error: pageManager.error
    }),
    [pageManager.currentPage, pageManager.handleChangePage, pageManager.isLoading, pageManager.error]
  );

  const pages = useMemo(
    () => ({
      [PageType.HOMEPAGE]: <HomePage />,

      [PageType.CHAT]: <MultiAI />,
      [PageType.SMART_BOOST]: <SmartBoost />,
      [PageType.ANALYTICS]: <SmartBoost />,

      // Boosters
      [PageType.FPS_BOOST]: <FpsBoostPage />,
      [PageType.CONNECTION]: <ConnectionPage />,
      [PageType.PRECISION]: <PrecisionPage />,
      [PageType.GAMES]: <GamesPage />,
      [PageType.FLUSHER]: <FlusherPage />,

      [PageType.PROFILE]: <ProfilePage />
    }),
    []
  );

  // Validação de autenticação
  if (authError) {
    return (
      <div className="flex items-center justify-center w-screen h-screen bg-zinc-950 text-white">
        <div className="text-center">
          <h2 className="text-xl font-semibold mb-2">{t('error.authErrorTitle')}</h2>
          <p className="text-red-400">{authError}</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="flex items-center justify-center w-screen h-screen bg-zinc-950 text-white">
        <LoadingState></LoadingState>
      </div>
    );
  }

  // Estado de erro geral
  if (pageManager.error) {
    return (
      <div className="flex items-center justify-center w-screen h-screen bg-zinc-950 text-white">
        <div className="text-center">
          <h2 className="text-xl font-semibold mb-2">{t('error.genericTitle')}</h2>
          <p className="text-red-400 mb-4">{pageManager.error}</p>
          <button
            onClick={() => pageManager.handleChangePage(PageType.HOMEPAGE)}
            className="bg-blue-600 hover:bg-blue-700 px-4 py-2 rounded"
          >
            {t('error.backToHome')}
          </button>
        </div>
      </div>
    );
  }

  return (
    <DashboardErrorBoundary>
      <UserProvider.Provider value={userProviderValues}>
        <PagesProvider.Provider value={pagesContextValue}>
          <section
            className="flex flex-row w-screen h-screen bg-zinc-950 text-white overflow-hidden"
            role="main"
            aria-label={t('ariaLabel')}
          >
            <DashboardHeader />

            <section className="flex bg-black flex-row items-start w-full h-full">
              {pageManager.isLoading ? (
                <div className="flex items-center justify-center w-full h-full">
                  <PageListingLoading />
                </div>
              ) : (
                <>
                  {Object.entries(pages).map(([pageType, component]) => (
                    <PageWrapper key={pageType} page={pageType as PageType} currentPage={pageManager.currentPage}>
                      {component}
                    </PageWrapper>
                  ))}
                </>
              )}
            </section>
          </section>
        </PagesProvider.Provider>
      </UserProvider.Provider>
    </DashboardErrorBoundary>
  );
}

// Hook para usar o contexto de páginas com validação
export const usePagesContext = () => {
  const context = useContext(PagesProvider);

  if (!context) {
    throw new Error('usePagesContext deve ser usado dentro de um PagesProvider');
  }

  return context;
};

// Constantes exportadas para uso externo
export const PAGE_TYPES = PageType;
export default DashboardPages;
