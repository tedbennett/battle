export type Board = {
  squares: number[][];
};

export type BoardMetadata = {
  colors: Record<number, string>;
};

export type Message = {
  type: number;
  board?: number[];
  colors?: Record<number, string>;
  diffs?: Diff[];
};

export type Diff = {
  row: number;
  col: number;
  team: number;
};
