/**
 * Code generators — run via `npm run gen:feature` / `gen:component` / `gen:hook`.
 * The point: nobody should ever hand-write the boilerplate for a new feature
 * slice, shared component, or shared hook. Ask Plop instead.
 */
export default function (plop) {
  plop.setGenerator('feature', {
    description: 'Scaffold a new feature slice under src/features/<name>',
    prompts: [
      {
        type: 'input',
        name: 'name',
        message: 'Feature name (e.g. "collections", "run", "socket-debugger"):',
      },
    ],
    actions: [
      {
        type: 'addMany',
        destination: 'src/features/{{dashCase name}}',
        base: 'plop-templates/feature',
        templateFiles: 'plop-templates/feature/**/*.hbs',
      },
    ],
  })

  plop.setGenerator('component', {
    description: 'Scaffold a shared component under src/components/ui/<Name>',
    prompts: [
      {
        type: 'input',
        name: 'name',
        message: 'Component name (e.g. "Button"):',
      },
    ],
    actions: [
      {
        type: 'add',
        path: 'src/components/ui/{{pascalCase name}}.tsx',
        templateFile: 'plop-templates/component/Component.tsx.hbs',
      },
      {
        type: 'add',
        path: 'src/components/ui/{{pascalCase name}}.test.tsx',
        templateFile: 'plop-templates/component/Component.test.tsx.hbs',
      },
    ],
  })

  plop.setGenerator('hook', {
    description: 'Scaffold a shared hook under src/hooks/use<Name>',
    prompts: [
      {
        type: 'input',
        name: 'name',
        message: 'Hook name, without the "use" prefix (e.g. "debounce"):',
      },
    ],
    actions: [
      {
        type: 'add',
        path: 'src/hooks/use{{pascalCase name}}.ts',
        templateFile: 'plop-templates/hook/useHook.ts.hbs',
      },
    ],
  })
}
