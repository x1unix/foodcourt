import { HttpParams } from '@angular/common/http';
import { isNil } from 'lodash';

const PARAM_KEYS = {
  orderBy: 'order_by',
  orderDir: 'order_dir',
  searchQuery: 'q',
  limit: 'limit',
  offset: 'offset'
};

export interface SearchQuery {
  orderBy?: string;
  orderDir?: string;
  searchQuery?: string;
  limit?: number;
  offset?: number;
}

export const queryFillParams = (q: SearchQuery) => {
  let params = new HttpParams();

  Object.keys(q).forEach((key: string) => {
    const queryParamName = PARAM_KEYS[key];

    if (!isNil(queryParamName)) {
      params = params.set(queryParamName, q[key]);
    }
  });

  return params;
};
