const filterFunc = (key: string, value: any) =>
  key === "__typename" ? undefined : value;

export default function omitTypename<T>(payload: T): T {
  return JSON.parse(JSON.stringify(payload), filterFunc);
}
