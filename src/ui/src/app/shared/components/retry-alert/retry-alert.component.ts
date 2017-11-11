import {Component, OnInit, ViewEncapsulation, Input, Output, EventEmitter} from '@angular/core';

/**
 * Request error alert
 */
@Component({
  selector: 'app-retry-alert',
  templateUrl: './retry-alert.component.html',
  styleUrls: ['./retry-alert.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class RetryAlertComponent implements OnInit {
  /**
   * Common text message
   * @type {string}
   */
  @Input() text = 'Failed to perform request';

  /**
   * Error details
   * @type {string}
   */
  @Input() errorMessage = 'Internal Server Error';

  /**
   * Retry click event
   * @type {EventEmitter<any>}
   */
  @Output() retry = new EventEmitter();

  constructor() { }

  ngOnInit() {
  }

}
