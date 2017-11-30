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
