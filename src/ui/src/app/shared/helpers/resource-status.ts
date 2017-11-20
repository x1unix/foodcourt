/**
 * Load status of external resource
 */

import {LoadStatus} from './load-status';

export class ResourceStatus {
  /**
   * Error message container
   *
   * @type {string}
   * @memberof LoadStatusComponent
   */
  public error: string = null;

  /**
   * Raw load status ID (used from LoadStatus enum)
   *
   * @type {number}
   * @memberof LoadStatusComponent
   */
  public status: number = LoadStatus.IDLE;


  /**
   * If request is loading
   *
   * @readonly
   * @type {boolean}
   * @memberof LoadStatusComponent
   */
  public get isLoading(): boolean {
    return this.status === LoadStatus.LOADING;
  }

  public set isLoading(val: boolean) {
    this.status = val ? LoadStatus.LOADING : LoadStatus.IDLE;
  }

  /**
   * If request is successfully loaded
   *
   * @readonly
   * @type {boolean}
   * @memberof LoadStatusComponent
   */
  public get isLoaded(): boolean {
    return this.status === LoadStatus.LOADED;
  }

  public set isLoaded(val: boolean) {
    this.status = val ? LoadStatus.LOADED : LoadStatus.FAILED;
  }

  /**
   *
   *
   * @readonly
   * @type {boolean}
   * @memberof LoadStatusComponent
   */
  public get isIdle(): boolean {
    return this.status === LoadStatus.IDLE;
  }

  public set isIdle(val: boolean) {
    this.status = val ? LoadStatus.IDLE : LoadStatus.LOADING;
  }

  /**
   * If request was failed
   *
   * @readonly
   * @type {boolean}
   * @memberof LoadStatusComponent
   */
  public get isFailed(): boolean {
    return this.status === LoadStatus.FAILED;
  }

  public set isFailed(val: boolean) {
    this.status = val ? LoadStatus.FAILED : LoadStatus.LOADED;
  }

  /**
   * List of load status ID's
   *
   * @readonly
   * @memberof LoadStatusComponent
   */
  public get STATUS() {
    return LoadStatus;
  }
}
