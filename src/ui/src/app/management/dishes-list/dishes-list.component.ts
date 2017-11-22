import {Component, OnInit, Input, OnDestroy, Output, EventEmitter} from '@angular/core';
import {IDish} from '../../shared/interfaces/dish';
import {DishesService} from '../services/dishes.service';
import {FormControl} from '@angular/forms';
import {Subscription} from 'rxjs/Rx';
import {DropEvent} from '../../shared/interfaces/drop-event';

const SEARCH_INPUT_TIMEOUT = 300;

@Component({
  selector: 'app-dishes-list',
  templateUrl: './dishes-list.component.html',
  styleUrls: ['./dishes-list.component.scss']
})
export class DishesListComponent implements OnInit, OnDestroy {

  @Input() items: IDish[] = [];

  @Input() selectedIds: number[] = [];

  @Input() catalogDropZone: string[] = [];

  @Input() cartDropZone: string[] = [];

  @Input() showDragBanner = false;

  @Input() allowDrag = true;

  @Output() drop = new EventEmitter<IDish>();

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

  onDropEnter(data: DropEvent<IDish>) {
    this.showDragBanner = false;
    const dish = data.dragData;
    this.drop.emit(dish);
  }

}
