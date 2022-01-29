# The Feature/Detail of This Controller

## Environment

「CRD も含めて何かカスタムコントローラを作る」ということを目的とし、当コントローラを作成しました。<br>
全く実用的なものではないですが、実装していく上で様々な知見を得ることができたと思います。<br>
ロジックなど処理が煩雑な箇所が多々ありますが、他のカスタムコントローラを作成する際に、当コントローラの実装が参考になれば（役に立てば）幸いです。

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

[main.go](../main.go) にて、"現在の時刻" と "次の日の 0:00" の差を取り、時間が来たら `dailyprocess.DeleteAttendanceBook()` を実行することで、定期実行を実現しています。<br>
`dailyprocess.DeleteAttendanceBook()` では、`DeleteAllOf()` を使用し、指定した Namespace 内の全ての ab リソースを削除する処理を実行しています。

## Future Feature 

* Leader election
* Read/Write custom ConfigMap
* Put logs to the other place (e.g. S3)