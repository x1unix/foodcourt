import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs/Observable';
import {ISettings} from './interfaces/settings';

const settingsUrl = '/api/settings';

@Injectable()
export class SettingsService {

  constructor(private http: HttpClient) { }

  /**
   * Gets settings
   * @returns {Observable<ISettings>}
   */
  getSettings(): Observable<ISettings> {
    return <Observable<ISettings>> this.http.get(settingsUrl);
  }

  /**
   * Save settings
   * @param {ISettings} settings
   * @returns {Observable<Object>}
   */
  saveSettings(settings: ISettings) {
    return this.http.post(settingsUrl, settings);
  }

}
