import * as moment from 'moment';

const SERVED_DATE_FMT = 'YYYYMMDD';
const DISPLAYED_DATE_FMT = 'DD/MM/YYYY';

export class DateTime {
  private _origin: moment.Moment;

  served: string;
  displayed: string;
  date: Date;

  constructor(m: moment.Moment) {
    this.origin = m;
  }

  get origin(): moment.Moment {
    return this._origin;
  }

  set origin(newVal: moment.Moment) {
    this._origin = newVal;
    this.displayed = newVal.format(DISPLAYED_DATE_FMT);
    this.served = newVal.format(SERVED_DATE_FMT);
    this.date = new Date(newVal.valueOf());
  }

  updateFromDate(date: Date) {
    this.origin = moment(date);
  }
}
