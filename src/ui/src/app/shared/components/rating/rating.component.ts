import {Component, OnInit, ViewEncapsulation, Input, EventEmitter, Output} from '@angular/core';

const RATING_MAX = 5;

interface Star {
  rate: number;
  shape: string;
  match: boolean;
}

/**
 * Stars rating component
 *
 * @example <app-rating [rating]="5" [readonly]="false" (change)="onRatingChange($value)"></app-rating>
 */
@Component({
  selector: 'app-rating',
  templateUrl: './rating.component.html',
  styleUrls: ['./rating.component.scss']
})
export class RatingComponent implements OnInit {

  /**
   * Initial rating
   * @type {number}
   */
  @Input() rating = 0;

  private currentValue = 0;

  /**
   * Selected rating value
   * @returns {number}
   */
  get value() {
    return this.currentValue;
  }

  /**
   * Selected rating value
   * @param {number} newVal
   */
  set value(newVal: number) {
    this.currentValue = newVal;
    this.recalcStarsIndex();
  }

  /**
   * Rating value change event
   * @type {EventEmitter<number>}
   */
  @Output() change = new EventEmitter<number>();

  /**
   * Max rating (stars count) (default: 5)
   * @type {number}
   */
  @Input() max = 5;

  /**
   * Is control at read-only state
   * @type {boolean}
   */
  @Input() readonly = false;

  /**
   * Formatted rating
   * @type {string}
   */
  prettyRating = '';

  starsIndex: Star[] = [];

  constructor() { }

  ngOnInit() {
    let rate = this.rating;

    if (rate > this.max) {
      rate = this.max;
    }

    this.preciseRating(rate);
    this.value = rate;
  }

  private preciseRating(rate: number) {
    this.prettyRating = (rate > 0) ? String(parseFloat('' + rate).toPrecision(2)) : 'N/A';
  }

  /**
   * Recalculate stars index
   */
  private recalcStarsIndex() {
    const newStars = [];

    for (let c = 1; c <= this.max; c++) {
      newStars.push({
        rate: c,
        shape: this.getButtonShape(c),
        match: this.value >= c
      });
    }

    if (this.rating === 0) {
      this.preciseRating(this.currentValue);
    }

    this.starsIndex = newStars;
  }

  getButtonShape(btnRating: number): string {
    const prefRate = btnRating - 1;

    return (this.value > prefRate) && (this.value < btnRating) ? 'half-star' : 'star';
  }

  /**
   * Set rating value
   * @param {number} newVal
   */
  setRating(newVal: number) {
    if (this.readonly) {
      return;
    }

    this.value = newVal;
    this.change.emit(newVal);
  }

}
