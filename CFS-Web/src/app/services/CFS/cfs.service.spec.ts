import { TestBed } from '@angular/core/testing';

import { CFSService } from './cfs.service';

describe('CFSService', () => {
  let service: CFSService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(CFSService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
