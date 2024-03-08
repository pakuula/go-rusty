// Copyright 2024 Nikolay Pakulin (@pakuula). All rights reserved.
// Use of this source code is governed by LGPL-3.0 licence.
// The text of the licence can be found in the LICENSE.txt file.

package os_r

import (
	"io/fs"
	"os"
	"time"

	"github.com/pakuula/go-rusty/option"
	"github.com/pakuula/go-rusty/result"
)

func Chdir(dir string) result.ResultVoid { return result.Void(os.Chdir(dir)) }
func Chmod(name string, mode os.FileMode) result.ResultVoid {
	return result.Void(os.Chmod(name, mode))
}
func Chown(name string, uid, gid int) result.ResultVoid {
	return result.Void(os.Chown(name, uid, gid))
}
func Chtimes(name string, atime time.Time, mtime time.Time) result.ResultVoid {
	return result.Void(os.Chtimes(name, atime, mtime))
}
func Create(fname string) result.Result[*os.File] {
	return result.Wrap(os.Create(fname))
}
func Executable() result.Result[string] { return result.Wrap(os.Executable()) }

func Getgroups() result.Result[[]int] { return result.Wrap(os.Getgroups()) }

func Getwd() result.Result[string]    { return result.Wrap(os.Getwd()) }
func Hostname() result.Result[string] { return result.Wrap(os.Hostname()) }
func Lchown(name string, uid, gid int) result.ResultVoid {
	return result.Void(os.Lchown(name, uid, gid))
}
func Link(oldname, newname string) result.ResultVoid { return result.Void(os.Link(oldname, newname)) }
func LookupEnv(key string) option.Option[string]     { return option.WrapOk(os.LookupEnv(key)) }
func Mkdir(name string, perm os.FileMode) result.ResultVoid {
	return result.Void(os.Mkdir(name, perm))
}
func MkdirAll(path string, perm os.FileMode) result.ResultVoid {
	return result.Void(os.MkdirAll(path, perm))
}
func MkdirTemp(dir, pattern string) result.Result[string] {
	return result.Wrap(os.MkdirTemp(dir, pattern))
}
func Open(fname string) result.Result[*os.File] {
	return result.Wrap(os.Open(fname))
}
func OpenFile(fname string, flag int, perm fs.FileMode) result.Result[*os.File] {
	return result.Wrap(os.OpenFile(fname, flag, perm))
}

func ReadFile(name string) result.Result[[]byte] { return result.Wrap(os.ReadFile(name)) }
func Readlink(name string) result.Result[string] { return result.Wrap(os.Readlink(name)) }
func Remove(name string) result.ResultVoid       { return result.Void(os.Remove(name)) }
func RemoveAll(path string) result.ResultVoid    { return result.Void(os.RemoveAll(path)) }
func Rename(oldpath, newpath string) result.ResultVoid {
	return result.Void(os.Rename(oldpath, newpath))
}
func Setenv(key, value string) result.ResultVoid { return result.Void(os.Setenv(key, value)) }
func Symlink(oldname, newname string) result.ResultVoid {
	return result.Void(os.Symlink(oldname, newname))
}
func Truncate(name string, size int64) result.ResultVoid {
	return result.Void(os.Truncate(name, size))
}
func Unsetenv(key string) result.ResultVoid { return result.Void(os.Unsetenv(key)) }
func UserCacheDir() result.Result[string]   { return result.Wrap(os.UserCacheDir()) }
func UserConfigDir() result.Result[string]  { return result.Wrap(os.UserConfigDir()) }
func UserHomeDir() result.Result[string]    { return result.Wrap(os.UserHomeDir()) }
func WriteFile(name string, data []byte, perm os.FileMode) result.ResultVoid {
	return result.Void(os.WriteFile(name, data, perm))
}
