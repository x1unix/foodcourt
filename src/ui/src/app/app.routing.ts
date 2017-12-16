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
import {OrderEditorComponent} from './container/order-editor/order-editor.component';
import {TodayComponent} from './container/today/today.component';
import {UsersManagerComponent} from './management/users-manager/users-manager.component';
import {SettingsComponent} from './management/settings/settings.component';
import {ProfileEditorComponent} from './auth/profile-editor/profile-editor.component';
import {ReportsComponent} from './reports/reports.component';


const APP_ROUTES: Routes = [
  {
    path: '',
    component: ContainerComponent,
    canActivate: [LoggedInGuard],
    children: [
      {
        path: '',
        component: DashboardComponent,
        data: {
          title: 'Orders'
        },
        children: [
          {
            path: '',
            component: TodayComponent,
            data: {
              title: 'Order for today'
            }
          },
          {
            path: 'orders',
            component: OrderEditorComponent,
            data: {
              title: 'Manage orders'
            }
          },
          {
            path: 'orders/:date/:userId',
            component: OrderEditorComponent,
            data: {
              title: 'Manage orders'
            }
          },
          {
            path: 'profile',
            component: ProfileEditorComponent,
            data: {
              title: 'Edit profile'
            }
          }
        ]
      },
      {
        path: 'management',
        component: ManagementComponent,
        canActivate: [AdminGuard],
        children: [
          {
            path: '',
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
          },
          {
            path: 'users',
            component: UsersManagerComponent,
            data: {
              title: 'Manage users'
            }
          },
          {
            path: 'system',
            component: SettingsComponent,
            data: {
              title: 'System settings'
            }
          }
        ]
      },
      {
        path: 'reports',
        component: ReportsComponent,
        canActivate: [AdminGuard]
      }
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
