import {Component, OnInit, Input, OnDestroy} from '@angular/core';
import {IDish} from '../../shared/interfaces/dish';
import {DishesService} from '../services/dishes.service';
import {FormControl} from '@angular/forms';
import {Subscription} from 'rxjs/Rx';
import * as lunr from 'lunr';

@Component({
  selector: 'app-dishes-list',
  templateUrl: './dishes-list.component.html',
  styleUrls: ['./dishes-list.component.scss']
})
export class DishesListComponent implements OnInit, OnDestroy {

  @Input() items: IDish[] = [];

  displayedItems: IDish[] = [];

  lunr: any = null;

  searchBoxControl: FormControl = null;

  searchChange$: Subscription;

  constructor(private dishesService: DishesService) { }

  ngOnInit() {
    this.displayedItems = this.items;

    this.searchBoxControl = new FormControl('');

    this.genSearchIndex();

    this.searchChange$ = this.searchBoxControl.valueChanges
      .debounceTime(500)
      .subscribe(newValue => this.search(newValue));
  }

  ngOnDestroy() {
    this.searchChange$.unsubscribe();
  }

  private genSearchIndex() {
    let docs = this.items;
    this.lunr = lunr(function () {
      this.ref('label');
      docs.forEach((item) => this.add(item));
    });
    docs = null;
  }

  /**
   * Gets category name
   * @param {number} catId Category id
   * @returns {string}
   */
  getDishCategory(catId: number) {
    return this.dishesService.getDishCategory(catId);
  }

  /**
   * Gets category color
   * @param {number} catId Category id
   * @returns {string}
   */
  getDishCategoryColor(catId: number) {
    return this.dishesService.getDishCategoryColor(catId);
  }

  search(searchQuery: string) {
    const oQuery = searchQuery.trim();

    if (oQuery.length === 0) {
      this.displayedItems = this.items;
      return;
    }

    this.displayedItems = this.lunr.search(oQuery);
  }

}
