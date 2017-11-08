import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {ContainerComponent} from './container/container.component';
import {LoggedInGuard} from './shared/guards/logged-in.guard';
import {AuthComponent} from './auth/auth.component';
import {DashboardComponent} from './container/dashboard/dashboard.component';

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
          title: 'Dashboard'
        }
      }
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
