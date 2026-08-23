import js from '@eslint/js'
import prettierConfig from 'eslint-config-prettier'
import importX from 'eslint-plugin-import-x'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import unusedImports from 'eslint-plugin-unused-imports'
import { globalIgnores } from 'eslint/config'
import globals from 'globals'
import tseslint from 'typescript-eslint'

export default tseslint.config(globalIgnores(['dist', 'wailsjs']), {
  files: ['**/*.{ts,tsx}'],
  extends: [
    js.configs.recommended,
    ...tseslint.configs.recommended,
    reactHooks.configs.flat['recommended-latest'],
    reactRefresh.configs.vite,
    importX.flatConfigs.recommended,
    importX.flatConfigs.typescript,
    // Must be last — turns off any rule that would fight Prettier's formatting.
    prettierConfig,
  ],
  languageOptions: {
    ecmaVersion: 2023,
    globals: globals.browser,
  },
  plugins: {
    'unused-imports': unusedImports,
  },
  settings: {
    'import-x/resolver': {
      typescript: true,
    },
  },
  rules: {
    // A finished feature never has an unused symbol lying around — this
    // catches leftovers from refactors that plain TS won't always flag.
    'unused-imports/no-unused-imports': 'error',

    // Forces the "@/..." alias (see tsconfig.app.json) instead of deep
    // relative paths — a feature module should never reach out via
    // ../../../ to grab something from another feature.
    'no-restricted-imports': [
      'error',
      {
        patterns: [
          {
            group: ['../*', '../../*', '../../../*'],
            message: 'Use the "@/..." alias instead of relative parent-directory imports.',
          },
        ],
      },
    ],

    'import-x/order': [
      'error',
      {
        groups: ['builtin', 'external', 'internal', 'parent', 'sibling', 'index'],
        pathGroups: [
          {
            pattern: '@/**',
            group: 'internal',
          },
        ],
        'newlines-between': 'always',
        alphabetize: { order: 'asc', caseInsensitive: true },
      },
    ],
  },
})
