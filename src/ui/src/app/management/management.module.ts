import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { SharedModule } from '../shared/shared.module';

import { ManagementComponent } from './management.component';
import { ItemsCatalogComponent } from './items-catalog/items-catalog.component';

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
    SharedModule
  ],
  declarations: [ManagementComponent, ItemsCatalogComponent]
})
export class ManagementModule { }
