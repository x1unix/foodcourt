import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';

/**
 * Alert message component
 *
 * @example <app-alert type="danger">Your message...</app-alert>
 */
@Component({
  selector: 'app-alert',
  templateUrl: './alert.component.html',
  styleUrls: ['./alert.component.scss']
})
export class AlertComponent implements OnInit {

  /**
   * Alert type
   * @type {string}
   */
  @Input() type = 'danger';

  /**
   * Alert icon
   * @type {string}
   */
  @Input() icon = 'exclamation-circle';

  /**
   * Show close button
   * @type {boolean}
   */
  @Input() dismissable = true;

  /**
   * Show action button
   * @type {boolean}
   */
  @Input() hasAction = false;

  /**
   * Action button test
   * @type {string}
   */
  @Input() actionText = 'Action';

  /**
   * Action button click event
   * @type {EventEmitter<any>}
   */
  @Output() actionClick = new EventEmitter();

  /**
   * Alert message close event
   * @type {EventEmitter<any>}
   */
  @Output() close = new EventEmitter();

  constructor() { }

  ngOnInit() {
  }

}
