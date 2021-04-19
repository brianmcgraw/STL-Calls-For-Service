import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { CFS, CFSResponse } from 'src/app/models/CFS/cfs';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class CFSService {

  constructor(private httpClient: HttpClient) { }

  public GetCFS(): Observable<CFS[]> {
    return this.httpClient.get<CFSResponse>("/api/cfs").pipe(
      map((cfsResponse: CFSResponse) => {
        // console.log(cfsResponse)
        return cfsResponse.cfs
      }),
      
    );
  }
  
}
