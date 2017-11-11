export enum DishType {
  soup,
  garnish,
  main,
  salad,
  special
}

export interface IDish {
  id: number;
  label: string;
  description: string;
  type: number;
  photoUrl: string;
}
