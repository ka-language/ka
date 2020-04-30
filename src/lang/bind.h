#ifdef __cplusplus
extern "C" {
#endif
  extern void Kill(void);
  extern char* Add(char*, char*, char*, int);
  extern char* Subtract(char*, char*, char*, int);
  extern char* Multiply(char*, char*, char*, int);
  extern char* Division(char*, char*, char*, int);
  extern char* Modulo(char*, char*, char*, int);
  extern char* Exponentiate(char*, char*, char*, int);
  extern char* Cactions(char*, char*);
  extern char* GetType(char*);
  extern char* ReturnInitC(char*);
  extern int IsLessC(char*, char*);
  extern char* CLex(char*);
  extern char* NQReplaceC(char*);
  extern int GetActNumC(char*);
  extern char* AddC(char*, char*);
  extern char* GetOp(char*);
  extern char* Similar(char*, char*, char*, char*, int, char*, char*);
  extern char* AddStrings(char*, char*, char*, int);
  extern char* SubtractStrings(char*, char*, char*, int);
  void bindCgo(char* actions, char* calc_params, char* dir);
#ifdef __cplusplus
}
#endif