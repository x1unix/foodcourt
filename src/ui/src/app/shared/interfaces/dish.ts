export enum DishType {
  soup,
  garnish,
  main,
  salad,
  special
}

export const DISH_TYPES = [
  'Soup',
  'Garnish',
  'Main',
  'Salad',
  'Special'
];

export const DISH_TYPE_COLORS = [
  'purple',
  'orange',
  'info',
  'success',
  'danger'
];

export const DISH_IMG_DEFAULT = '/assets/dish_no_image.jpg';

export interface IDish {
  id?: number;
  label: string;
  description: string;
  type: number;
  photoUrl: string;
}
