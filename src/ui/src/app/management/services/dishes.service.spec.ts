import { TestBed, inject } from '@angular/core/testing';

import { DishesService } from './dishes.service';

describe('DishesService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [DishesService]
    });
  });

  it('should be created', inject([DishesService], (service: DishesService) => {
    expect(service).toBeTruthy();
  }));
});
