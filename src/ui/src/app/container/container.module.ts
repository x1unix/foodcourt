import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { ClarityModule } from 'clarity-angular';
import { SharedModule } from '../shared/shared.module';

import { ContainerComponent } from './container.component';
import { HeaderComponent } from './header/header.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { OrderEditorComponent } from './order-editor/order-editor.component';
import { TodayComponent } from './today/today.component';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    RouterModule,
    ClarityModule,
    SharedModule
  ],
  declarations: [ContainerComponent, HeaderComponent, DashboardComponent, OrderEditorComponent, TodayComponent]
})
export class ContainerModule { }
