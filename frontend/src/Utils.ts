export class Utils {
  public static groupArray<T>(arr: T[], group: number): T[][] {
    let chunk = Math.ceil(arr.length / group)
    let groups = [], i = 0, n = arr.length;
    while (i < n) {
      groups.push(arr.slice(i, i += chunk));
    }
    return groups;
  }
}
