import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ManagementComponent } from './management.component';
import { ItemsCatalogComponent } from './items-catalog/items-catalog.component';

@NgModule({
  imports: [
    CommonModule
  ],
  declarations: [ManagementComponent, ItemsCatalogComponent]
})
export class ManagementModule { }
