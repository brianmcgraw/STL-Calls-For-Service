import { HttpClient } from '@angular/common/http';
import { Component, OnInit, QueryList, ViewChild, ViewChildren } from '@angular/core';
import { GoogleMap, MapInfoWindow, MapMarker}  from '@angular/google-maps';
import { Observable, of } from 'rxjs';
import { catchError, map, mergeMap, toArray } from 'rxjs/operators';
import { CFS } from './models/CFS/cfs';
import { CFSService } from './services/CFS/cfs.service';
import { environment } from '../environments/environment'

export interface MapsInfo {
  LatLng: google.maps.LatLngLiteral,
  Id: string
  Description: string
  EventTime: string
  Ward: string
  Neighborhood: string
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  apiLoaded: Observable<boolean>;
  @ViewChild(GoogleMap, { static: false }) map?: GoogleMap
  @ViewChildren(MapInfoWindow) infoWindowsView: QueryList<MapInfoWindow>;
  infoWindowOptions: google.maps.InfoWindowOptions = {
    disableAutoPan: true
  }

  public zoom = 12;
  public center: google.maps.LatLngLiteral = {
    lat: 38.62,
    lng: -90.243
  }
  // https://timdeschryver.dev/blog/google-maps-as-an-angular-component
  public options: google.maps.MapOptions = {
    zoomControl: true,
    scrollwheel: false,
    disableDoubleClickZoom: true,
    mapTypeId: 'roadmap',
    mapTypeControl: false,
    maxZoom: 15,
    minZoom: 8,
  }
  
  

  constructor(private cfsService: CFSService, httpClient: HttpClient) {
    this.apiLoaded = httpClient.jsonp(`https://maps.googleapis.com/maps/api/js?key=${environment.MAPS_KEY}`, 'callback')
        .pipe(
          map(() => true),
          catchError(() => of(false)),
        );
  }
  

  ngOnInit() {
    this.getCFS();
  }

public markers$: Observable<MapsInfo[]>;
  public getCFS() {
    this.markers$ = this.cfsService.GetCFS().pipe(
      mergeMap((cfs: CFS[]) => cfs),
      map(
        (cfs: CFS) => {
        let marker: google.maps.LatLngLiteral;
        marker = {
          lat: cfs.latitude,
          lng: cfs.longitude
        }
        return {
          LatLng: marker,
          
          Label: "label",
          Content: "<div>Here is some html</div>",
          Id: cfs.Id,
          Description: cfs.description,
          EventTime: String(cfs.eventTime),
          Ward: cfs.ward,
          Neighborhood: cfs.neighborhood
        }
      }),
      toArray()
      )
  }

  openInfoWindow(marker: MapMarker, windowIndex: number) {
    let curIdx = 0;
    this.infoWindowsView.forEach((window: MapInfoWindow) => {
      if (windowIndex === curIdx) {
        window.open(marker);
        curIdx++;
      } else {
        curIdx++;
      }
    });
  }

  closeInfoWindow(marker: MapMarker, windowIndex: number) {
    /// stores the current index in forEach
    
    let curIdx = 0;
    this.infoWindowsView.forEach((window: MapInfoWindow) => {
    
      if (windowIndex === curIdx) {
    
        window.close();
        curIdx++;
      } else {
        curIdx++;
      }
    });
  }
}
