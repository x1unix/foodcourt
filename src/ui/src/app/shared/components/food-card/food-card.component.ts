import {Component, OnInit, ViewEncapsulation, Input} from '@angular/core';

@Component({
  selector: 'app-food-card',
  templateUrl: './food-card.component.html',
  styleUrls: ['./food-card.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class FoodCardComponent implements OnInit {

  // Behavior props

  /**
   * Can element be selected
   * @type {boolean}
   */
  @Input() selectable = true;

  /**
   * Is element selected
   * @type {boolean}
   */
  @Input() selected = false;

  /**
   * Can element be rated
   * @type {boolean}
   */
  @Input() rateable = true;

  constructor() { }

  ngOnInit() {
  }

  toggleSelected() {
    this.selected = !this.selected;
  }

}
