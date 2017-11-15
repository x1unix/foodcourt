import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {ContainerComponent} from './container/container.component';
import {LoggedInGuard} from './shared/guards/logged-in.guard';
import {AuthComponent} from './auth/auth.component';
import {DashboardComponent} from './container/dashboard/dashboard.component';
import {ManagementComponent} from './management/management.component';
import {AdminGuard} from './shared/guards/admin.guard';
import {ItemsCatalogComponent} from './management/items-catalog/items-catalog.component';
import {MenuEditorComponent} from './management/menu-editor/menu-editor.component';

const APP_ROUTES: Routes = [
  {
    path: '',
    component: ContainerComponent,
    canActivate: [LoggedInGuard],
    children: [
      {
        path: 'dashboard',
        component: DashboardComponent,
        data: {
          title: 'Dashboard'
        }
      },
      {
        path: 'management',
        component: ManagementComponent,
        canActivate: [AdminGuard],
        children: [
          {
            path: 'schedule',
            component: MenuEditorComponent,
            data: {
              title: 'Manage menu'
            }
          },
          {
            path: 'catalog',
            component: ItemsCatalogComponent,
            data: {
              title: 'Manage dishes'
            }
          }
        ]
      },
      // {
      //   path: '**',
      //   redirectTo: '/dashboard'
      // }
    ]
  },
  {
    path: 'auth',
    component: AuthComponent,
    data: {
      title: 'Authorization'
    }
  },
  {
    path: '**',
    redirectTo: '/dashboard'
  }
];

@NgModule({
  imports: [RouterModule.forRoot(APP_ROUTES)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
