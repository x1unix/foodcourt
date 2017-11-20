import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { ClarityModule } from 'clarity-angular';
import { SharedModule } from '../shared/shared.module';

import { ContainerComponent } from './container.component';
import { HeaderComponent } from './header/header.component';
import { DashboardComponent } from './dashboard/dashboard.component';

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
    ClarityModule,
    SharedModule
  ],
  declarations: [ContainerComponent, HeaderComponent, DashboardComponent]
})
export class ContainerModule { }
