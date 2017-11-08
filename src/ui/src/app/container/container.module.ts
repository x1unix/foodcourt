import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

import { LoggedInGuard } from '../shared/guards/logged-in.guard';
import { ContainerComponent } from './container.component';

@NgModule({
  imports: [
    CommonModule,
    RouterModule
  ],
  declarations: [ContainerComponent]
})
export class ContainerModule { }
