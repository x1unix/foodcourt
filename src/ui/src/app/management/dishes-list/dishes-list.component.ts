import {Component, OnInit, Input, OnDestroy} from '@angular/core';
import {IDish} from '../../shared/interfaces/dish';
import {DishesService} from '../services/dishes.service';
import {FormControl} from '@angular/forms';
import {Subscription} from 'rxjs/Rx';

const SEARCH_INPUT_TIMEOUT = 300;

@Component({
  selector: 'app-dishes-list',
  templateUrl: './dishes-list.component.html',
  styleUrls: ['./dishes-list.component.scss']
})
export class DishesListComponent implements OnInit, OnDestroy {

  @Input() items: IDish[] = [];

  @Input() selectedIds: number[] = [];

  displayedItems: IDish[] = [];

  searchBoxControl: FormControl = null;

  searchChange$: Subscription;

  constructor(private dishesService: DishesService) { }

  ngOnInit() {
    this.displayedItems = this.items;

    this.searchBoxControl = new FormControl('');

    this.searchChange$ = this.searchBoxControl.valueChanges
      .debounceTime(SEARCH_INPUT_TIMEOUT)
      .subscribe(newValue => this.search(newValue));
  }

  ngOnDestroy() {
    this.searchChange$.unsubscribe();
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
    let oQuery = searchQuery.trim();

    if (oQuery.length === 0) {
      this.displayedItems = this.items;
      return;
    }

    let searchExpr = new RegExp(oQuery, 'gi');
    this.displayedItems = this.items.filter((i) => searchExpr.test(i.label));

    searchExpr = oQuery = null;
  }

}
