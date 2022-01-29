# attendance-book-in-kubernetes

> いつもどおりのある日の事 君は突然立ち上がり言った<br>
> 「Kubernetes で勤怠管理をしよう」

## About this controller

「Kubernetes 上で勤怠管理アプリを動かすのではなく、"Kubernetes" を使って勤怠管理を行うこと」を目標とするコントローラ。

メリット : 『俺たちは Kubernetes を使ってるぞ！』という優越感（？）に浸ることができる。<br>
デメリット :  全くいらない

当コントローラの機能・詳細は [こちら](docs/feature.md)

## How to deploy

1. Kubernetes クラスターを用意する
2. 下記コマンドで当コントローラをデプロイする
```
$ git clone https://github.com/a2ush/attendance-book-in-kubernetes.git
$ cd attendance-book-in-kubernetes
$ make docker-build docker-push IMG=<registry>/<project-name>:tag
$ make deploy IMG=<registry>/<project-name>:tag
```

## How to use

### Employee

従業員は Custome Resource である `AttendanceBook (短縮名 ab)` をデプロイし、当日の出欠状況を伝える。 
```
$ cat sample-manifest.yaml
apiVersion: office.a2ush.dev/v1alpha1
kind: AttendanceBook
metadata:
  name: sample-user
spec:
  attendance: absent      # Input present or absent
  reason: "Feel sleepy"   # Optional
  
$ kubectl apply -f sample-manifest.yaml
```
present, absent 以外の内容を入力した場合は、無効な状況として出欠報告が受け付けられない。
```
$ kubectl apply -f invalid-attendance.yaml
The AttendanceBook "test-user3" is invalid: spec.attendance: Unsupported value: "almostpresent": supported values: "present", "absent"
```

また、雇用者が報告を受け付けている Namespace（デフォルト: default） 以外で `AttendanceBook` をデプロイした場合は、当コントローラにより即時削除されるので気をつける。
```
$ cat not-permitted.yaml
apiVersion: office.a2ush.dev/v1alpha1
kind: AttendanceBook
metadata:
  name: test-user4
  namespace: notpermitted
spec:
  attendance: present

$ kubectl apply -f not-permitted.yaml
attendancebook.office.a2ush.dev/test-user4 created

$ kubectl get ab -n notpermitted
No resources found in notpermitted namespace.
$ kubectl get event -n notpermitted
LAST SEEN   TYPE     REASON    OBJECT                      MESSAGE
78s         Normal   Deleted   attendancebook/test-user4   Deleted resource notpermitted/test-user4 due to the namespace that is not allowed to deploy.
```

もし、最初に伝えた内容と異なる状況になった場合（例: 病気による早退、出勤可能になったなど）は、従業員はすぐに自身の `AttendanceBook` を更新する。（更新方法は問わない）
```
Ex)
$ kubectl patch ab sample-user --type='json' \
  -p='[{"op": "replace", "path": "/spec/attendance", "value": "present"}, {"op": "replace", "path": "/spec/reason", "value": "My eyes are opened completely"}]'
```

デフォルトでは日本時間の 0:00 に全ての `AttendanceBook` が削除されるため、従業員は日毎に出欠状況を報告する必要がある。

### Employer

雇用者は `kubectl get ab` コマンドを実行し、各従業員の出勤状況を確認する。
```
$ kubectl get ab
NAME          ATTENDANCE   REASON
sample-user   present      My eyes are opened completely
test-user1    present      Work is my life
test-user2    absent       BLANK
```

そして、雇用者は `kubectl describe ab` コマンドを実行することで、従業員がステータス変更を行ったかどうかを確認することができる。
```
$ kubectl describe ab sample-user
...
Events:
  Type    Reason   Age    From                       Message
  ----    ------   ----   ----                       -------
  Normal  Created  2m46s  attendancebook-controller  Created resource. Attendance/Reason: absent/Feel sleepy
  Normal  Updated  118s   attendancebook-controller  Updated resource. Attendance/Reason: present/My eyes are opened completely
```

### Optional

 - 環境変数を設定し、デフォルトの挙動を変えることができる [[manifest]](config/manager/manager.yaml)
```
        env:
        - name: SPECIFIED_NAMESPACE
          value: "default"
          name: TIMEZONE
          value: "Asia/Tokyo"
```
`SPECIFIED_NAMESPACE` で指定した Namespace 以外に ab リソースをデプロイした場合、そのリソースは当コントローラによって即時削除される。<br>
`TIMEZONE` は、どのタイムゾーンの 0:00 を指しているかを示すものである。
