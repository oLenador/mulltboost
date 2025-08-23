// // src/presentation/features/boosters/providers/batch-manager.provider.tsx
// 
// import React, { createContext, useContext, useRef, useEffect, ReactNode } from 'react';
// import { Provider as JotaiProvider } from 'jotai';
// import { BoosterBatchManager } from '../domain/booster-batch.manager';
// 
// interface BoosterBatchManagerContextType {
//   manager: BoosterBatchManager | null;
// }
// 
// const BoosterBatchManagerContext = createContext<BoosterBatchManagerContextType>({
//   manager: null
// });
// 
// interface BoosterBatchManagerProviderProps {
//   children: ReactNode;
// }
// 
// export const BoosterBatchManagerProvider: React.FC<BoosterBatchManagerProviderProps> = ({ children }) => {
//   const managerRef = useRef<BoosterBatchManager | null>(new BoosterBatchManager());
// 
//   useEffect(() => {
//     // Initialize manager if not already created
//     if (!managerRef.current) {
//       managerRef.current = new BoosterBatchManager();
//       
//       // Sync with backend on initialization
//       managerRef.current.syncWithBackend().catch(console.error);
//       
//       // Set up periodic sync (every 30 seconds)
//       const syncInterval = setInterval(() => {
//         managerRef.current?.syncWithBackend().catch(console.error);
//       }, 30000);
// 
//       return () => {
//         clearInterval(syncInterval);
//       };
//     }
//   }, []);
// 
//   return (
//     <JotaiProvider>
//       <BoosterBatchManagerContext.Provider value={{ manager: managerRef.current }}>
//         {children}
//       </BoosterBatchManagerContext.Provider>
//     </JotaiProvider>
//   );
// };
// 
// export const useBoosterBatchManagerContext = () => {
//   const context = useContext(BoosterBatchManagerContext);
//   if (!context) {
//     throw new Error('useBoosterBatchManagerContext must be used within BoosterBatchManagerProvider');
//   }
//   return context;
// };