function getRandomInt(max: number) {
  return Math.floor(Math.random() * max);
}

function randomPickAndRemove<T>(arr: T[]): {value: T, arr: T[]} {
  const index = getRandomInt(arr.length);
  return {
    value: arr[index],
    arr: arr.slice(0, index).concat(arr.slice(index+1, arr.length))
  }
}

//valueがnum個ある配列をランダムな並び順で生成する
export function genRandomOrderArray<T>(elements: {value: T, num: number}[]): T[] {
  const alignmentValuesArray = elements.reduce((pre: T[], current: {value: T, num: number}) => {
    const arr = Array<T>(current.num).fill(current.value);
    return pre.concat(arr);
  }, []);

  return Array(alignmentValuesArray.length).fill(null).reduce((pre: {reaming: T[], randomArray: T[]}) => {
    const {value, arr} = randomPickAndRemove(pre.reaming);
    return {
      reaming: arr,
      randomArray: pre.randomArray.concat([value]),
    }
  }, {reaming: alignmentValuesArray, randomArray: []}).randomArray;
}
