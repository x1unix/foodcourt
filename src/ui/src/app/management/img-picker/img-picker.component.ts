import { Component, OnInit, Input, Output, EventEmitter, ViewChild } from '@angular/core';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';

@Component({
  selector: 'app-img-picker',
  templateUrl: './img-picker.component.html',
  styleUrls: ['./img-picker.component.scss']
})
export class ImgPickerComponent extends LoadStatusComponent implements OnInit {
  @Input() src: string = null;

  @Output() upload = new EventEmitter<string>();

  @ViewChild('fileInput') fileInput: HTMLInputElement;

  constructor() {
    super();
  }

  ngOnInit() {}

  onFilePick() {
    if (this.fileInput.files.length === 0) {
      return;
    }

    const formData = new FormData();
    const file = this.fileInput.files[0];

    formData.append('image', file);

    this.isLoading = true;
  }

}
