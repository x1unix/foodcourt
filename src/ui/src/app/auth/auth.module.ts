import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AuthComponent } from './auth.component';
import { SharedModule } from '../shared/shared.module';
import { ClarityModule } from 'clarity-angular';
import { FormsModule } from '@angular/forms';

@NgModule({
  imports: [
    CommonModule,
    SharedModule,
    FormsModule,
    ClarityModule
  ],
  declarations: [AuthComponent]
})
export class AuthModule { }
