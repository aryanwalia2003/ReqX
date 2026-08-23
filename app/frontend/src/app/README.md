# App shell

Root-level wiring only: providers (theme, query client, etc. as they get
added), the top-level layout, and the app's entry composition. If you're
adding a feature, it almost certainly belongs in `src/features/`, not here.
