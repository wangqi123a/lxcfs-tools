// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// iSulad-lxcfs-toolkit is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.
// Description: exec prestart mount hook
// Author: zhangsong
// Create: 2019-01-18

// go base main package
package main

import (
	"fmt"
	"io/ioutil"
	"isula.org/isulad-lxcfs-toolkit/libmount"
	"os"
	"strconv"

	isulad_lxcfs_log "github.com/sirupsen/logrus"
)

func prestartMountHook(pid int, rootfs string) error {
	lxcfssubpath, err := ioutil.ReadDir("/var/lib/lxc/lxcfs/proc")
	if err != nil {
		isulad_lxcfs_log.Errorf("Prase lxcfs dir failed: %v", err)
		return err
	}

	initMountns, err := os.Readlink("/proc/1/ns/mnt")
	if err != nil {
		return fmt.Errorf("read init mount namespace fail: %v", err)
	}
	mountns, err := os.Readlink("/proc/" + strconv.Itoa(pid) + "/ns/mnt")
	if err != nil {
		return fmt.Errorf("read container mount namespace fail: %v", err)
	}
	if initMountns == mountns {
		return fmt.Errorf("container pid changed: container mount namespace is same as init namespace")
	}

	var valuePaths []string
	var valueMountPaths []string
	for _, value := range lxcfssubpath {
		valuePaths = append(valuePaths, fmt.Sprintf("%s/proc/%s", rootfs, value.Name()))
		valueMountPaths = append(valueMountPaths, fmt.Sprintf("/var/lib/lxc/lxcfs/proc/%s", value.Name()))
	}

	if err := libmount.NsExecMount(strconv.Itoa(pid), valueMountPaths, valuePaths); err != nil {
		isulad_lxcfs_log.Errorf("mount %v into container error: %v", valueMountPaths, err)
		return err
	}
	return nil
}
