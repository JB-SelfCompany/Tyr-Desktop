/**
 * useI18n - Hook that integrates react-i18next with config store
 * Automatically syncs language preference with backend configuration
 */

import { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useConfigStore, type Language } from '../store/configStore';

/**
 * Hook that provides i18n functionality with config sync
 * Automatically synchronizes language changes with backend config
 *
 * @example
 * ```tsx
 * const { t, language, changeLanguage } = useI18n();
 *
 * return (
 *   <div>
 *     <h1>{t('welcome')}</h1>
 *     <button onClick={() => changeLanguage('ru')}>
 *       Switch to Russian
 *     </button>
 *   </div>
 * );
 * ```
 */
export function useI18n() {
  const { t, i18n } = useTranslation();
  const config = useConfigStore((state) => state.config);
  const setLanguage = useConfigStore((state) => state.setLanguage);

  // Sync i18n with config on mount and config changes
  useEffect(() => {
    if (config && config.language !== i18n.language) {
      i18n.changeLanguage(config.language);
    }
  }, [config, i18n]);

  // Update HTML lang attribute when language changes
  useEffect(() => {
    const htmlElement = document.documentElement;
    if (htmlElement) {
      htmlElement.setAttribute('lang', i18n.language);
    }
  }, [i18n.language]);

  /**
   * Change language and save to backend config
   */
  const changeLanguage = async (language: Language) => {
    try {
      // Update i18n immediately for responsive UI
      await i18n.changeLanguage(language);
      // Save to backend config
      await setLanguage(language);
    } catch (error) {
      console.error('Failed to change language:', error);
      throw error;
    }
  };

  return {
    t,
    i18n,
    language: (config?.language || i18n.language) as Language,
    changeLanguage,
  };
}

/**
 * Hook that returns only the translation function
 * Useful when you only need translations without language management
 *
 * @example
 * ```tsx
 * const t = useTranslate();
 * return <h1>{t('common.title')}</h1>;
 * ```
 */
export function useTranslate() {
  const { t } = useTranslation();
  return t;
}

/**
 * Hook that returns current language
 *
 * @example
 * ```tsx
 * const currentLanguage = useCurrentLanguage();
 * ```
 */
export function useCurrentLanguage(): Language {
  const config = useConfigStore((state) => state.config);
  const { i18n } = useTranslation();
  return (config?.language || i18n.language) as Language;
}
