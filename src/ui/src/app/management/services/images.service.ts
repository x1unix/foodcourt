import { Injectable } from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Observable} from 'rxjs/Observable';
import {IMessage} from '../../shared/interfaces/message';

@Injectable()
export class ImagesService {

  constructor(private http: HttpClient) { }

  upload(file: File) {
    const formData = new FormData();
    formData.append('image', file);

    const headers = new HttpHeaders();
    headers.set('Content-Type', 'multipart/form-data');

    return this.http.post('/api/media', formData);
  }

}
