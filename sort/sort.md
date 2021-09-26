#### 冒泡排序
```c
void BubbleSort(SqList *L)
{
    int i,j;
    for (i=1; i<L->length; i++) {
        for (j=L->length-1; j > i; j--) {
            if (L->r[j] < L->r[j-1]) {
                tmp = L->r[j-1];
                L->r[j-1] = L->r[j];
                L->r[j] = tmp;
            }
        }
    }
}
```

#### 插入排序
```c
void SelectSort(SqList *L)
{
    int i,j,min;
    for (i=1; i<L->length; i++) {
        min = i;
        for (j=i+1; j<= L->length; j++) {
            if (L->r[j] < L->[min]) {
                min = j;
            }
        }
        
        if (i!=min) {
            tmp = L->r[min];
            L->r[min] = L->r[i];
            L->r[i] = tmp;
        }
        
    }
}
```

#### 直接插入排序
 - 思想是 如果这个数小，单独拿出来，大的往前移动(这不是还是冒泡吗?)
```c
void InsertSort(SqList *L)
{
    int i,j;
    for (i=2; i<=L->length; i++) {
        if (L->r[i] < L->r[i-1]) {
            L->r[0] = L->r[i];
            for (j=i-1; L->r[j]>L->r[0]; j--) {
                L->[j+1] = L->[j];
            }
            L->r[j+1] = L->r[0];
        }
    }
}
```

#### 希尔排序


#### 堆排序
 - 堆是完全二叉树:
    + 大顶堆: 每个结点的值都大于或等于其左右孩子结点的值
    + 小顶堆: 每个结点的值都小于或等于其左右孩子结点的值称为小顶堆
    
 - 根结点一定是堆中所有结点最大(小)者

