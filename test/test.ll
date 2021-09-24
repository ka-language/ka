%tusk.string = type { i8*, i32 }
%tuskclass.Main = type { void (%tuskclass.Main)* }

@.Main_d = global void (%tuskclass.Main)* null
@.Main_main = global void ()* null

define %tusk.string @tusk.newstring(i8* %sptr, i32 %slen) {
0:
	%1 = alloca %tusk.string
	%2 = getelementptr %tusk.string, %tusk.string* %1, i32 0, i32 0
	store i8* %sptr, i8** %2
	%3 = getelementptr %tusk.string, %tusk.string* %1, i32 0, i32 1
	store i32 %slen, i32* %3
	%4 = load %tusk.string, %tusk.string* %1
	ret %tusk.string %4
}

define void @_tusk_init() {
0:
	store void ()* @tv_5, void ()** @.Main_main
	ret void
}

define %tuskclass.Main @tuskclass.new.Main() {
0:
	%1 = alloca %tuskclass.Main
	store void (%tuskclass.Main)* @tv_4, void (%tuskclass.Main)** @.Main_d
	%2 = load %tuskclass.Main, %tuskclass.Main* %1
	ret %tuskclass.Main %2
}

declare void @tv_1(%tuskclass.Main %0)

declare void @tv_2()

declare void @tv_3(%tuskclass.Main %0)

define void @tv_4(%tuskclass.Main %a) {
0:
	ret void
}

define void @tv_5() {
0:
	%1 = call %tuskclass.Main @tuskclass.new.Main()
	%2 = alloca %tuskclass.Main
	store %tuskclass.Main %1, %tuskclass.Main* %2
	%3 = getelementptr %tuskclass.Main, %tuskclass.Main* %2, i32 0, i32 0
	%4 = load void (%tuskclass.Main)*, void (%tuskclass.Main)** %3
	%5 = load %tuskclass.Main, %tuskclass.Main* %2
	call void %4(%tuskclass.Main %5)
	ret void
}

define void @main() {
0:
	call void @_tusk_init()
	ret void
}
