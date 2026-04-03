/**
 * Normalize an array to return an empty array if it's not an array.
 * @param arr - The array to normalize.
 * @returns An empty array if arr is not an array, arr otherwise.
 */
export const NormalizeArray = (arr: unknown) => {
    return Array.isArray(arr) ? arr : [];
}