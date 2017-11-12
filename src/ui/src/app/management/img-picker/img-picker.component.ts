import { Component, OnInit, Input, Output, EventEmitter, ViewChild, ElementRef } from '@angular/core';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import { ImagesService } from '../services/images.service';
import {IMessage} from '../../shared/interfaces/message';
import {WebHelperService} from '../../shared/services/web-helper.service';

const MAX_FILE_SIZE = 2 * 1024 * 1024; // 2 Mb

@Component({
  selector: 'app-img-picker',
  templateUrl: './img-picker.component.html',
  styleUrls: ['./img-picker.component.scss'],
  providers: [
    ImagesService
  ]
})
export class ImgPickerComponent extends LoadStatusComponent implements OnInit {
  @Input() src: string = null;

  @Output() uploadSuccess = new EventEmitter<string>();

  @Output() uploadFailed = new EventEmitter<string>();

  @Output() uploadStart = new EventEmitter();

  @ViewChild('fileInput') fileInput: ElementRef;

  constructor(private images: ImagesService, private web: WebHelperService) {
    super();
  }

  ngOnInit() {}

  onFilePick() {
    const input: HTMLInputElement = this.fileInput.nativeElement;

    if (input.files.length === 0) {
      return;
    }

    const file: File = input.files[0];

    if (file.size > MAX_FILE_SIZE) {
      window.alert(`Max allowed file size is 2MB`);
      return;
    }

    this.uploadImage(file);
  }

  private uploadImage(file: File) {
    this.isLoading = true;
    this.uploadStart.emit();
    this.images.upload(file).subscribe(
      (data: IMessage) => {
        this.uploadSuccess.emit(data.msg);
        this.src = data.msg;
        this.isLoaded = true;
      },
      (err) => {
        this.isFailed = true;
        const msg = this.web.extractResponseError(err);
        alert(`Failed to upload an image: ${msg}`);
        this.uploadFailed.emit(msg);
      }
    );
  }

}
