package dailyprocess

import (
	"context"
	"log"
	"os"

	"github.com/a2ush/attendance-book-in-kubernetes/controllers"
	"sigs.k8s.io/controller-runtime/pkg/client"

	officev1alpha1 "github.com/a2ush/attendance-book-in-kubernetes/api/v1alpha1"
)

func DeleteAttendanceBook(ctx context.Context, cli client.Client) error {

	specified_namespace := controllers.GetNamespace()
	err := cli.DeleteAllOf(ctx, &officev1alpha1.AttendanceBook{}, client.InNamespace(specified_namespace))
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Deleted all AttendanceBooks in %s.\n", specified_namespace)
	return nil
}

func GetTimezone() string {
	timezone, found := os.LookupEnv("TIMEZONE")
	if !found {
		timezone = "Asia/Tokyo"
	}
	return timezone
}
