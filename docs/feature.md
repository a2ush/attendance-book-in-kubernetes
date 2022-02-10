# The Feature/Detail of This Controller

## Environment

「CRD も含めて何かカスタムコントローラを作る」ということを目的とし、当コントローラを作成しました。<br>
ロジック等の処理が煩雑な箇所が散見されるかもしれませんが、他のカスタムコントローラを作成する際に、当コントローラの実装が参考になれば幸いです。

このコントローラは `Kubebuilder` で作成されたもので、テスト環境は EKS となっています。

## Option
 - 環境変数を設定し、デフォルトの挙動を変えることができます [[manifest]](../config/manager/manager.yaml)
```
        env:
        - name: SPECIFIED_NAMESPACE
          value: "default"
          name: TIMEZONE
          value: "Asia/Tokyo"
```
`SPECIFIED_NAMESPACE` で指定した Namespace 以外に ab リソースをデプロイした場合、そのリソースは当コントローラによって即時削除されます。<br>
`TIMEZONE` は、どのタイムゾーンの 0:00 を指しているかを示すものです。

## Work

### Reconcile の処理

Custome Resource `AttendanceBook(短縮名 ab)` を検知し、controller.go 内の処理によって Status の更新や Event の作成を行っています。

### 定期実行の処理（日毎に ab リソースを削除する処理）

[main.go](../main.go) にて、"現在の時刻" と "次の日の 0:00" の差を取り、時間が来たら `dailyprocess.DeleteAttendanceBook()` を実行することで、定期実行を実現しています。
goroutine による並列処理を行っているため、Reconcile の処理は妨げられません。<br>
`dailyprocess.DeleteAttendanceBook()` では、`DeleteAllOf()` を使用し、指定した Namespace 内の全ての ab リソースを削除する処理を実行しています。
なお、`DeleteAllOf()` を実行するためには、Kubernetes RBAC の `deletecollection` verb の実行を許可する必要があります。


### 従業員リストに記載の無い従業員の AttendanceBook を削除する処理

従業員リストは ConfigMap `employee-list` であり、当コントローラの Deployment リソースにマウントされています。 <br>
従業員名は `employeeList` slice に格納され、ab リソースが最初に作成されたとき（desire.Status.Attendance が "" のとき）に、作成された ab リソース名が slice の要素にあるかどうかを確認しています。

### リーダー選出

Kuberbuilder v3 では、デフォルトでリーダー選出の機能が有効になっています。<br>
ConfigMap `094dacd8.a2ush.dev` に現在のリーダーの情報が記載されます。
```
$ k get po -n attendance-book-in-kubernetes-system 
NAME                                                              READY   STATUS    RESTARTS   AGE
attendance-book-in-kubernetes-controller-manager-7cc588799m7t4k   1/1     Running   0          77s
attendance-book-in-kubernetes-controller-manager-7cc588799zvqzj   1/1     Running   0          77s

$ kubectl get cm -n attendance-book-in-kubernetes-system 094dacd8.a2ush.dev -oyaml | grep leader
    control-plane.alpha.kubernetes.io/leader: '{"holderIdentity":"attendance-book-in-kubernetes-controller-manager-7cc588799zvqzj_a053b184-2db7-4827-95b4-3e079dc19e2d","leaseDurationSeconds":15,"acquireTime":"2022-02-10T12:58:42Z","renewTime":"2022-02-10T14:25:21Z","leaderTransitions":0}
```
リーダーに選ばれたコントローラは、以下のようなログが記録されます。
```
$ kubectl logs -n attendance-book-in-kubernetes-system attendance-book-in-kubernetes-controller-manager-7cc588799zvqzj
...
I0210 12:58:42.782951       1 leaderelection.go:248] attempting to acquire leader lease attendance-book-in-kubernetes-system/094dacd8.a2ush.dev...
I0210 12:58:42.809579       1 leaderelection.go:258] successfully acquired lease attendance-book-in-kubernetes-system/094dacd8.a2ush.dev
...
```

## Future Feature 

* Put logs to the other place (e.g. Amazon S3 bucket)