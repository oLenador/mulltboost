// i18n.ts
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import { getDefaultStore } from 'jotai';
import { userDataAtom } from '@/core/store/user-data.store';
import dashboardEn from './en/dashboard.json';
import dashboardPt from './pt/dashboard.json';
import dashboardEs from './es/dashboard.json';
import dashboardRu from './ru/dashboard.json';
import dashboardPtBr from './pt-BR/dashboard.json';

import homepageEn from './en/homepage.json';
import homepagePt from './pt/homepage.json';
import homepageEs from './es/homepage.json';
import homepageRu from './ru/homepage.json';
import homepagePtBr from './pt-BR/homepage.json';

import boostersEn from './en/boosters.json';
import boostersPt from './pt/boosters.json';
import boostersEs from './es/boosters.json';
import boostersRu from './ru/boosters.json';
import boostersPtBr from './pt-BR/boosters.json';

// Pega valor inicial do Jotai
const store = getDefaultStore();
const { language } = store.get(userDataAtom);

i18n
  .use(initReactI18next)
  .init({
    lng: language,
    fallbackLng: 'en',
    interpolation: { escapeValue: false },
    resources: {
      en: { dashboard: dashboardEn, homepage: homepageEn, boosters: boostersEn},
      pt: { dashboard: dashboardPt, homepage: homepagePt, boosters: boostersPt},
      "pt-BR": { dashboard: dashboardPtBr, homepage: homepagePtBr, boosters: boostersPtBr},
      ru: { dashboard: dashboardRu, homepage: homepageRu, boosters: boostersRu},
      es: { dashboard: dashboardEs, homepage: homepageEs, boosters: boostersEs},
    },
  });

store.sub(userDataAtom, () => {
  const { language } = store.get(userDataAtom);
  i18n.changeLanguage(language);
});

export default i18n;
