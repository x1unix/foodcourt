import { TestBed, inject } from '@angular/core/testing';

import { WebHelperService } from './web-helper.service';

describe('WebHelperService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [WebHelperService]
    });
  });

  it('should be created', inject([WebHelperService], (service: WebHelperService) => {
    expect(service).toBeTruthy();
  }));
});
