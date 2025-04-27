import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, map, Observable, of, Subject } from 'rxjs';
import { AttackInput } from "../interfaces/AttackInput";
import { AttackOutput } from '../interfaces/AttackOutput';

@Injectable({
  providedIn: 'root'
})
export class ApiClientService {

  constructor(private client : HttpClient) { }

  private baseBackendUrl = 'http://localhost:8080/';

  private apiSubject$ = new Subject<AttackOutput>();

  getFakeApiResponse() : Observable<AttackOutput> {
    return this.apiSubject$.asObservable();
  }

  getFakeApi(request: AttackInput) : Observable<boolean> {
    return this.client.get<AttackOutput>(this.baseBackendUrl + 'api/fake').pipe(
      map((response: AttackOutput) => {
        setTimeout(() => {
          this.apiSubject$.next(response);
        }, 2000); // Simulate a delay of 1 second
        return true;
      }),
      catchError((error) => {
        console.error('Error:', error);
        return of(false); // Return false or handle the error as needed
      })
    );
  }
}
