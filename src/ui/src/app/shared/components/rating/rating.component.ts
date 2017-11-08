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
 * @example <app-rating [value]="5" [readonly]="false" (change)="onRatingChange($value)"></app-rating>
 */
@Component({
  selector: 'app-rating',
  templateUrl: './rating.component.html',
  styleUrls: ['./rating.component.scss']
})
export class RatingComponent implements OnInit {

  private rating = 0;

  /**
   * Rating value
   * @returns {number}
   */
  get value() {
    return this.rating;
  }

  /**
   * Rating value
   * @param {number} newVal
   */
  @Input() set value(newVal: number) {
    this.rating = newVal;
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
  prettyRating = '0.0';

  starsIndex: Star[] = [];

  constructor() { }

  ngOnInit() {
    if (this.rating > this.max) {
      this.rating = this.max;
    }
  }

  /**
   * Recalculate stars index
   */
  private recalcStarsIndex() {
    const newStars = [];
    const rating = String(parseFloat('' + this.rating).toPrecision(2));

    for (let c = 1; c <= this.max; c++) {
      newStars.push({
        rate: c,
        shape: this.getButtonShape(c),
        match: this.rating >= c
      });
    }

    this.prettyRating = rating;
    this.starsIndex = newStars;
  }

  getButtonShape(btnRating: number): string {
    const prefRate = btnRating - 1;

    return (this.rating > prefRate) && (this.rating < btnRating) ? 'half-star' : 'star';
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
