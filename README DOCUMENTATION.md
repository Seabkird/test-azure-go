/src
├── components/          # Composants UI globaux, réutilisables et "agnostiques" (Boutons, Cards génériques)
│   └── Welcome/         # Exemple de structure "Co-location" (Code, Style, Test, Story au même endroit)
├── features/            # Logique métier découpée par domaine (Vertical Slicing)
│   ├── demo/            # Fonctionnalité de démo
│   └── users/           # Domaine "Utilisateurs"
│       └── UserList.tsx # Composant métier spécifique aux utilisateurs
├── lib/                 # Configuration d'infrastructure et utilitaires tiers
│   └── api.ts           # Client HTTP (Axios/Fetch), intercepteurs, gestion d'erreurs
├── pages/               # Les "Vues" ou "Pages" complètes (liées aux routes URL)
│   ├── Home.page.tsx    # Page d'accueil
│   └── Users.page.tsx   # Page de gestion des utilisateurs
├── test-utils/          # Utilitaires de configuration pour les tests (Vitest/Jest)
├── App.tsx              # Composant racine (Providers globaux, Layout principal)
├── main.tsx             # Point d'entrée de React (Mount dans le DOM)
├── Router.tsx           # Définition des routes de l'application
├── theme.ts             # Surcharge du thème Mantine (Couleurs, Typographie)
└── vite-env.d.ts        # Types TypeScript pour Vite

