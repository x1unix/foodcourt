import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import { ClarityModule } from 'clarity-angular';
import { SharedModule } from '../shared/shared.module';

import { ReportsComponent } from './reports.component';
import {ReportsService} from './reports.service';

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
    FormsModule,
    ReactiveFormsModule,
    ClarityModule,
    SharedModule
  ],
  declarations: [ReportsComponent],
  providers: [
    ReportsService
  ]
})
export class ReportsModule { }
