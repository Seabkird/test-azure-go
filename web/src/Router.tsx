import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { HomePage } from './pages/Home.page';
import { UsersPage } from './pages/Users.page';
import { MainLayout } from './components/layout/MainLayout';

const router = createBrowserRouter([
  {
    path: '/',
    element: <MainLayout />, // On applique le Layout ici (Parent)
    children: [
      {
        path: '/',
        element: <HomePage />,
      },
      {
        path: '/users',
        element: <UsersPage  />,
      },
    ],
  },
]);

export function Router() {
  return <RouterProvider router={router} />;
}
